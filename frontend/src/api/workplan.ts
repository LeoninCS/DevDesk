// src/api/workplan.ts
import http from "./http";

export interface TodoItem {
  id: number;
  content: string;
  done: boolean;
}

export interface WorkPlanNewResponse {
  hash: string;
}

// GET /api/workplan/new
export function createWorkPlan() {
  return http.get<WorkPlanNewResponse>("/workplan/new");
}

// GET /api/workplan/:hash
export function fetchWorkPlanTodos(hash: string) {
  return http.get<{ todos: TodoItem[] }>(`/workplan/${hash}`);
}

// POST /api/workplan/add
export function addWorkPlanTodo(hash: string, content: string) {
  return http.post("/workplan/add", { hash, content });
}

// POST /api/workplan/done
export function toggleWorkPlanTodo(hash: string, id: number) {
  return http.post("/workplan/done", { hash, id });
}

// POST /api/workplan/edit
export function editWorkPlanTodo(hash: string, id: number, content: string) {
  return http.post("/workplan/edit", { hash, id, content });
}

// POST /api/workplan/delete
export function deleteWorkPlanTodo(hash: string, id: number) {
  return http.post("/workplan/delete", { hash, id });
}
