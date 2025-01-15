package global

import "regexp"

// 高危命令检测
func FindDangerousCommandsPos(content string, dangerous_commands []string) ([][]int, []string) {
	var positions [][]int
	var matchedCommands []string

	for _, pattern := range dangerous_commands {
		re, err := regexp.Compile(pattern)
		if err != nil {
			// TODO: info remind
			continue
		}
		matches := re.FindAllStringIndex(content, -1)
		for _, match := range matches {
			start, end := match[0], match[1]-1
			positions = append(positions, []int{start, end})
			matchedCommands = append(matchedCommands, content[start:end+1])
		}
	}
	return positions, matchedCommands
}

var DangerousCommandsList = []string{
	`.*rm\s+-[r,f,rf].*`,
	`.*lvremove\s+-f.*`,
	`.*poweroff.*`,
	`.*shutdown\s+-[f,F,h,k,n,r,t,C].*`,
	`.*pvremove\s+-f.*`,
	`.*vgremove\s+-f.*`,
	`.*exportfs\s+-[a,u].*`,
	`.*umount.nfs+.*.+-[r,f,rf].*`,
	`.*mv+.*.+/dev/null.*`,
	`.*reboot.*`,
	`.*rmmod\s+-[a,s,v,f,w].*`,
	`.*dpkg-divert+.*.+-remove.*`,
	`.*dd.*`,
	`.*mkfs.*`,
	`.*vmo.*`,
	`.*init.*`,
	`.*halt.*`,
	`.*fasthalt.*`,
	`.*fastboot.*`,
	`.*startsrc.*`,
	`.*stopsrc.*`,
	`.*chkconfig.*`,
	`.*off.*`,
	`.*refresh.*`,
	`.*umount.*`,
	`.*rmdev.*`,
	`.*chdev.*`,
	`.*extendvg.*`,
	`.*reducevg.*`,
	`.*importvg.*`,
	`.*exportvg.*`,
	`.*mklv.*`,
	`.*rmlv.*`,
	`.*rmfs.*`,
	`.*chfs.*`,
	`.*installp.*`,
	`.*instfix.*`,
	`.*crontab.*`,
	`.*cfgmgr.*`,
	`.*mknod.*`,
}