
export interface Extention {
  name: string;
  permission: string;
  type: string;
  url: string;
  parentName?: string;
}

export type ExtArr = Array<Extention>;