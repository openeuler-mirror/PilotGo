// 历史版本清单
export interface HistoryItem {
  description: string;
  content: string;
  version: string;
  UpdatedAt: string;
  id: number; // 版本id
  scriptid?: number; // 脚本id
}
