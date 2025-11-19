// src/api/codeshare.ts
import http from "./http";

export interface UploadCodePayload {
  author: string;
  language: string;
  content: string;
  destroy_time: number;
}

export interface UploadCodeResponse {
  hash: string;
  url?: string;
}

export function uploadCode(data: UploadCodePayload) {
  // 实际请求：<baseURL>/codeshare/upload
  return http.post<UploadCodeResponse>("/codeshare/upload", data);
}

export function getCodeByHash(hash: string) {
  // 实际请求：<baseURL>/codeshare/code/:hash
  return http.get(`/codeshare/code/${hash}`);
}
