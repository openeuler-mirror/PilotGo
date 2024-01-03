
// 校验登录账号
export let checkAccount = (rule: any, value: any, callback: Function) => {
    if (!value) {
        return callback();
    }
    if (value === "admin") {
        return callback();
    }
    let patt = /^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$/;
    if (!patt.test(value)) {
        return callback(new Error("邮箱格式错误"))
    }
    return callback()
};
