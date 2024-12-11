import { NextResponse } from "next/server";
import { db } from "~/server/db";
import { feeds } from "~/server/db/schema";
import Parser from "rss-parser";
import { nanoid } from "nanoid";
import { auth } from "~/server/auth";
import { eq } from "drizzle-orm";

let parser = new Parser();

export const GET = auth(async function GET(req) {
  if (req.auth == null) return NextResponse.error();
  const res = await db
    .select()
    .from(feeds)
    .where(eq(feeds.user, req.auth.user.id));
  return Response.json(res);
});

export const POST = auth(async function POST(request) {
  if (request.auth == null) return NextResponse.error();

  const { url } = await request.json();
  let parsed = await parser.parseURL(url.toString());

  const feed: typeof feeds.$inferInsert = {
    nid: nanoid(),
    title: parsed.title ?? url,
    url,
    image: parsed.image?.url,
    summary: parsed.description,
    authors: parsed.items?.map((item) => item.author).join(","),
    user: request.auth.user.id,
  };
  await db.insert(feeds).values(feed);

  return NextResponse.json({ ok: true });
});
