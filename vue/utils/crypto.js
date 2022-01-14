import CryptoJS from 'crypto-js/crypto-js'
//msg 需要被对称加密的明文
//key aes 对称加密的密钥  必须是16长度,为了和后端交互 key字符串必须是16进制字符串,否在给golang进行string -> []byte带来困难
 export function encrypt(msg, key) {
    key = PaddingLeft(key, 16);//保证key的长度为16byte,进行'0'补位
    key = CryptoJS.enc.Utf8.parse(key);
    // 加密结果返回的是CipherParams object类型
    // key 和 iv 使用同一个值
    var encrypted = CryptoJS.AES.encrypt(msg, key, {
        iv: key,
        mode: CryptoJS.mode.CBC,// CBC算法
        padding: CryptoJS.pad.Pkcs7 //使用pkcs7 进行padding 后端需要注意
    });
    // ciphertext是密文,toString()内传编码格式,比如Base64,这里用了16进制
    // 如果密文要放在 url的参数中 建议进行 base64-url-encoding 和 hex encoding, 不建议使用base64 encoding
    return  encrypted.ciphertext.toString(CryptoJS.enc.Hex)  //后端必须进行相反操作

}
// 确保key的长度,使用 0 字符来补位
// length 建议 16 24 32
function PaddingLeft(key, length){
    let  pkey= key.toString();
    let l = pkey.length;
    if (l < length) {
        pkey = new Array(length - l + 1).join('0') + pkey;
    }else if (l > length){
        pkey = pkey.slice(length);
    }
    return pkey;
}