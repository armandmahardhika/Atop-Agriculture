import { BasicResponse } from "./types";
import { path } from "ramda";
export const isFine = (r: any): boolean => r.code === 0;

export const errorMessage = (o: any): string => {
  const p = ["payload", "reason"];
  const msg = path(p, o) as string;
  return msg ? msg : "Unknown error";
};
