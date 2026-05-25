import { useMutation, useQueryClient } from '@tanstack/react-query';

import { acceptConsent, changeVote, reportSurvey, verifyPin, vote } from '@/services/voting';

type VoteVariables = {
  optionId: string;
  access_token?: string;
};

type ChangeVoteVariables = {
  newOptionId: string;
  access_token?: string;
};

type VerifyPinVariables = {
  pin: string;
};

type ReportSurveyVariables = {
  reason: string;
  details?: string;
  access_token?: string;
};

export function useVote(surveyId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ optionId, access_token }: VoteVariables) =>
      vote(surveyId, optionId, access_token),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ['survey', surveyId] });
      void queryClient.invalidateQueries({ queryKey: ['survey-results', surveyId] });
      void queryClient.invalidateQueries({ queryKey: ['feed'] });
    },
  });
}

export function useChangeVote(surveyId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ newOptionId, access_token }: ChangeVoteVariables) =>
      changeVote(surveyId, newOptionId, access_token),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ['survey', surveyId] });
      void queryClient.invalidateQueries({ queryKey: ['survey-results', surveyId] });
      void queryClient.invalidateQueries({ queryKey: ['feed'] });
    },
  });
}

export function useVerifyPin(surveyId: string) {
  return useMutation({
    mutationFn: ({ pin }: VerifyPinVariables) => verifyPin(surveyId, pin),
  });
}

export function useAcceptConsent() {
  return useMutation({
    mutationFn: () => acceptConsent(),
  });
}

export function useReportSurvey(surveyId: string) {
  return useMutation({
    mutationFn: ({ reason, details, access_token }: ReportSurveyVariables) =>
      reportSurvey(surveyId, reason, details, access_token),
  });
}
