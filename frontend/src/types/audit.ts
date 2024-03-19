export interface AuditItem {
  CreatedAt: string;
  DeletedAt: string;
  ID: number;
  /* 
  * @params Isempty 
  * 0:has children logs  
  * 1:no children logs
  *  */
  Isempty: number; 
  hasChildren?: boolean;
  UpdatedAt: string;
  action: string;
  log_uuid: string;
  message: string;
  module: string;
  parent_uuid: string;
  status: string;
  user_id: number;
}