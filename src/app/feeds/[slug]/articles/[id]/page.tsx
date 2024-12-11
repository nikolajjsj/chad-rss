export default async function ArticlePage({
  params,
}: {
  params: Promise<{ slug: string, id: string }>
}) {
  const slug = (await params).slug
  const id = (await params).id

  return <div>My Article: {id} from Feed: {slug}</div>
}
