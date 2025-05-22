import api from '../utils/axiosSetup';

export interface RequestLog {
  id: number;
  created_at: string;
  ip: string;
  method: string;
  path: string;
  status: number;
  user_agent: string;
  category: string;
  details: string;
  response_time: number;
}

export async function fetchRecentRequestLogs(limit = 50, category?: string): Promise<RequestLog[]> {
  const params: any = { limit };
  if (category) params.category = category;
  const res = await api.get('/api/analytics/requests/recent', { params });
  return res.data;
}
