// The IOTServer basic response
export interface BasicResponse {
  version: string;
  code: number;
  status: string;
  command: string;
  payload?: object;
}

export interface JSONObject {
  [k: string]: string | number;
}
