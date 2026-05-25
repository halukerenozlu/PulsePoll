import { api } from './api';
import type { ConsentResponse, PinVerifyResponse, ReportResponse, VoteResponse } from '@/types';

function authHeaders(access_token?: string): HeadersInit | undefined {
  return access_token ? { Authorization: `Bearer ${access_token}` } : undefined;
}

export function vote(
  surveyId: string,
  optionId: string,
  access_token?: string,
): Promise<VoteResponse> {
  return api<VoteResponse>(`/api/v1/surveys/${surveyId}/vote`, {
    method: 'POST',
    headers: authHeaders(access_token),
    body: JSON.stringify({ option_id: optionId }),
    credentials: 'include',
  });
}

export function changeVote(
  surveyId: string,
  newOptionId: string,
  access_token?: string,
): Promise<VoteResponse> {
  return api<VoteResponse>(`/api/v1/surveys/${surveyId}/vote`, {
    method: 'PUT',
    headers: authHeaders(access_token),
    body: JSON.stringify({ new_option_id: newOptionId }),
    credentials: 'include',
  });
}

export function verifyPin(surveyId: string, pin: string): Promise<PinVerifyResponse> {
  return api<PinVerifyResponse>(`/api/v1/surveys/${surveyId}/pin/verify`, {
    method: 'POST',
    body: JSON.stringify({ pin }),
    credentials: 'include',
  });
}

export function acceptConsent(): Promise<ConsentResponse> {
  return api<ConsentResponse>('/api/v1/consent/accept', {
    method: 'POST',
    credentials: 'include',
  });
}

export function reportSurvey(
  surveyId: string,
  reason: string,
  details?: string,
  access_token?: string,
): Promise<ReportResponse> {
  return api<ReportResponse>(`/api/v1/surveys/${surveyId}/report`, {
    method: 'POST',
    headers: authHeaders(access_token),
    body: JSON.stringify({ reason, ...(details ? { details } : {}) }),
    credentials: 'include',
  });
}
