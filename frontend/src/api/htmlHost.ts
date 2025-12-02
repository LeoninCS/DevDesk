import http from "./http";

export interface UploadHtmlResponse {
  url: string;
  full_url?: string;
  filename?: string;
  limit_bytes?: number;
}

export function uploadHtmlFile(file: File) {
  const formData = new FormData();
  formData.append("file", file);

  return http.post<UploadHtmlResponse>("/html/upload", formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
}

