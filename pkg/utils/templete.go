package utils

import (
	"reflect"
	"sort"
	"strings"

	yaml "gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

const DefaultTagName = "comment"
const (
	PilotGoDisabledComment CommentsTagFlag = "disabled_comment" //关闭注释
	PilotGoHeadComment     CommentsTagFlag = "head_comment"     //注释在上
	PilotGoLineComment     CommentsTagFlag = "line_comment"     //注释在同一行
	PilotGoFootComment     CommentsTagFlag = "foot_comment"     //注释在下
)

type YamlOpeartor struct {
	structInfo interface{}
	options    *Options
}

type CommentsTagFlag string

func (ctf CommentsTagFlag) enabled() bool {
	return ctf != PilotGoDisabledComment
}

type Options struct {
	CommentsTagFlag CommentsTagFlag
	OmitEmpty       bool
	DefaultTagName  string
}

func newOptions(opts ...Option) *Options {
	opt := &Options{
		CommentsTagFlag: PilotGoDisabledComment,
		OmitEmpty:       true,
		DefaultTagName:  DefaultTagName,
	}

	for _, f := range opts {
		f(opt)
	}

	return opt
}

type Option func(*Options)

func WithCommentsTagFlag(ctf CommentsTagFlag) Option {
	return func(o *Options) {
		o.CommentsTagFlag = ctf
	}
}
func WithDefaultTagName(defaultTagName string) Option {
	return func(o *Options) {
		if defaultTagName != "" {
			o.DefaultTagName = defaultTagName
		}
	}
}
func WithOmitEmpty(value bool) Option {
	return func(o *Options) {
		o.OmitEmpty = value
	}
}

func NewYamlOpeartor(config interface{}, opts ...Option) *YamlOpeartor {
	return &YamlOpeartor{
		structInfo: config,
		options:    newOptions(opts...),
	}
}

func (e *YamlOpeartor) ConvertYamlNode() (*yaml.Node, error) {
	node, err := convertYamlNode(e.structInfo, e.options)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (e *YamlOpeartor) Marshal() ([]byte, error) {
	if e.options.CommentsTagFlag == PilotGoDisabledComment {
		return yaml.Marshal(e.structInfo)
	}

	node, err := e.ConvertYamlNode()
	if err != nil {
		return nil, err
	}
	if node.Kind == yaml.MappingNode && len(node.Content) == 0 && node.FootComment != "" && e.options.CommentsTagFlag.enabled() {
		res := ""
		if node.HeadComment != "" {
			res += node.HeadComment + "\n"
		}
		lines := strings.Split(res+node.FootComment, "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
				continue
			}
			lines[i] = "# " + line
		}
		return []byte(strings.Join(lines, "\n")), nil
	}

	return yaml.Marshal(node)
}

func isEmpty(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}
	switch value.Kind() {
	case reflect.Ptr:
		return value.IsNil()
	case reflect.Map:
		return len(value.MapKeys()) == 0
	case reflect.Slice:
		return value.Len() == 0
	default:
		return value.IsZero()
	}
}

func isNil(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}
	switch value.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
func convertYamlNode(in interface{}, options *Options) (*yaml.Node, error) {
	node := &yaml.Node{}
	if n, ok := in.(*yaml.Node); ok {
		return n, nil
	}
	if m, ok := in.(yaml.Marshaler); ok && !isNil(reflect.ValueOf(in)) {
		res, err := m.MarshalYAML()
		if err != nil {
			return nil, err
		}
		if n, ok := res.(*yaml.Node); ok {
			return n, nil
		}
		in = res
	}
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		node.Kind = yaml.MappingNode
		k := v.Type()
		for i := 0; i < v.NumField(); i++ {
			if !v.Field(i).CanInterface() {
				continue
			}
			comment := k.Field(i).Tag.Get(options.DefaultTagName)
			tag := k.Field(i).Tag.Get("yaml")
			tagParts := strings.Split(tag, ",")
			fieldName := tagParts[0]
			tagParts = tagParts[1:]
			if fieldName == "" {
				fieldName = strings.ToLower(k.Field(i).Name)
			}
			if fieldName == "-" {
				continue
			}
			isEmpty := isEmpty(v.Field(i))
			isSkip := false
			isInline := false
			flow := false
			for _, part := range tagParts {
				if part == "omitempty" && isEmpty && options.OmitEmpty {
					isSkip = true
				}
				if part == "inline" {
					isInline = true
				}

				if part == "flow" {
					flow = true
				}
			}
			var value interface{}
			if v.Field(i).CanInterface() {
				value = v.Field(i).Interface()
			}
			if isSkip {
				continue
			}
			var yamlStyle yaml.Style
			if flow {
				yamlStyle |= yaml.FlowStyle
			}
			if isInline {
				child, err := convertYamlNode(value, options)
				if err != nil {
					return nil, err
				}
				if child.Kind == yaml.MappingNode || child.Kind == yaml.SequenceNode {
					appendYamlNodes(node, child.Content...)
				}
			} else if err := mergeComments(node, fieldName, value, comment, yamlStyle, options); err != nil {
				return nil, err
			}
		}
	case reflect.Map:
		node.Kind = yaml.MappingNode
		keys := v.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})
		for _, k := range keys {
			element := v.MapIndex(k)
			value := element.Interface()

			if err := mergeComments(node, k.Interface(), value, "", 0, options); err != nil {
				return nil, err
			}
		}
	case reflect.Slice:
		node.Kind = yaml.SequenceNode
		nodes := make([]*yaml.Node, v.Len())
		for i := 0; i < v.Len(); i++ {
			element := v.Index(i)
			var err error
			nodes[i], err = convertYamlNode(element.Interface(), options)
			if err != nil {
				return nil, err
			}
		}
		appendYamlNodes(node, nodes...)
	default:
		if err := node.Encode(in); err != nil {
			return nil, err
		}
	}
	return node, nil
}

func appendYamlNodes(node *yaml.Node, nodes ...*yaml.Node) {
	if node.Content == nil {
		node.Content = []*yaml.Node{}
	}
	node.Content = append(node.Content, nodes...)
}

func mergeComments(dest *yaml.Node, fieldName, in interface{}, fieldComent string, style yaml.Style, options *Options) error {
	yamlNode, err := convertYamlNode(fieldName, options)
	if err != nil {
		return err
	}
	value, err := convertYamlNode(in, options)
	if err != nil {
		return err
	}
	value.Style = style
	if options.CommentsTagFlag.enabled() {
		addCommentPosition(yamlNode, fieldComent, options.CommentsTagFlag)
	}
	appendYamlNodes(dest, yamlNode, value)
	return nil
}

func addCommentPosition(node *yaml.Node, commentValue string, tagFlag CommentsTagFlag) {
	if tagFlag.enabled() {
		switch tagFlag {
		case PilotGoDisabledComment:
			klog.Error("disabled comment tag!")
		case PilotGoHeadComment:
			node.HeadComment = commentValue
		case PilotGoLineComment:
			node.LineComment = commentValue
		case PilotGoFootComment:
			node.FootComment = commentValue
		default:
			klog.Warning("no support comment tag!")
		}

	}
}
