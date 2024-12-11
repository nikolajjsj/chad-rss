import { useMutation, useQuery } from "@tanstack/react-query";
import { getQueryClient } from "~/app/providers";
import { feeds } from "~/server/db/schema";

export const useFeeds = () => {
  return useQuery({
    queryKey: ['feeds'],
    queryFn: async () => {
      const response = await fetch('/api/feeds')
      const data = await response.json()
      return data as typeof feeds.$inferSelect[];
    }
  })
}

export const useCreateFeed = () => {
  const queryClient = getQueryClient()

  return useMutation({
    mutationKey: ['createFeed'],
    mutationFn: async ({ url }: { url: string }) => {
      const response = await fetch('/api/feeds', {
        method: 'POST',
        body: JSON.stringify({ url }),
      })
      return await response.json()
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({
        queryKey: ['feeds']
      })
    }
  })
}
