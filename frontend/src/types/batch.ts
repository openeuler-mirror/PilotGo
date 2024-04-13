// 批次列表
export interface BatchItem {
  CreatedAt:String;
  DeletedAt:null;
  Depart:String;
  DepartName: String;
  ID: number;
  UpdatedAt: String;
  description: String;
  manager: String;
  name: String;
}

// 批次详情列表
export interface BatchMachineInfo {
  CPU: String;
  departid: number;
  id:number;
  ip: String;
  machineuuid: String;
  maintstatus: String;
  runstatus: String;
  sysinfo: String;
}