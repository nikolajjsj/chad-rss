'use client'

import { useFeeds } from "~/hooks/api/feeds";

export default function FeedsPage() {
  const { data, isLoading } = useFeeds();

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return (
    <main className="flex flex-col justify-center items-center">
      {data?.map((feed) => (
        <div key={feed.id} className="flex items-center space-x-4">
          {feed.image ? (
            <img src={feed.image} alt={feed.title} className="w-12 h-12" />
          ) : null}
          <div>
            <h2 className="text-lg font-semibold">{feed.title}</h2>
          </div>
        </div>
      ))}
    </main>
  );
}
