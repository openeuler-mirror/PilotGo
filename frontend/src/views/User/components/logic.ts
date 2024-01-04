
// 校验邮箱
export let checkEmail = (rule: any, value: any, callback: Function) => {
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


// 校验手机号
export let checkPhone = (rule: any, value: any, callback: Function) => {
    if (!value) {
        return callback();
    }

    let reg = /^[1]([3-9])[0-9]{9}$/;
    if (!reg.test(value)) {
        return callback(new Error("手机号格式错误"))
    }
    return callback()
};