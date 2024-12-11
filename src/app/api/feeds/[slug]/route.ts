import { eq } from "drizzle-orm";
import { NextRequest } from "next/server";
import { db } from "~/server/db";
import { feeds } from "~/server/db/schema";

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ slug: string }> },
) {
  const slug = (await params).slug;

  const feed = await db
    .select()
    .from(feeds)
    .where(eq(feeds.nid, slug))
    .limit(1);

  return Response.json(feed);
}
