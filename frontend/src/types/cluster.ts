export interface MachineInfo {
  cpu: string;
  departid: number;
  departname: string;
  id: number;
  ip: string;
  maintstatus: string;
  runstatus: string;
  systeminfo: string;
  uuid: string;
}

export interface DeptTree {
  id: number;
  pid: number;
  label: string;
  children?: DeptTree[]
}