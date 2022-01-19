// 校验ip
  export let checkIP = (rule, value, callback, obj) => {
    let reg = /^([1-9]|[1-9]\d|1\d{2}|2[0-1]\d|22[0-3])(\.(\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])){3}$/;
    regex(rule, value, callback, obj, reg);
  };

// 校验邮箱
export let checkEmail = (rule, value, callback) => {
  if (!value) {
    return callback();
  }
  let patt = /^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$/;
  regex(rule, value, callback, {}, patt);
};
  
  function regex(rule, value, callback, obj, reg) {
    if (!value) {
      callback();
    }
    if (!reg.test(value)) {
      return callback(new Error());
    }
    callback();
  }