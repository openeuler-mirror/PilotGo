/* 格式化Date */
export function formatDate(date: Date) {
  const year = date.getFullYear(); // 获取年份
  const month = String(date.getMonth() + 1).padStart(2, "0"); // 获取月份，注意要加1，并补零
  const day = String(date.getDate()).padStart(2, "0"); // 获取日期，并补零
  const hours = String(date.getHours()).padStart(2, "0"); // 获取小时，并补零
  const minutes = String(date.getMinutes()).padStart(2, "0"); // 获取分钟，并补零
  const seconds = String(date.getSeconds()).padStart(2, "0"); // 获取秒，并补零

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`; // 拼接成所需格式
}
