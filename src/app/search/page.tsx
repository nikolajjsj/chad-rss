import Form from 'next/form'
import { Button } from '~/components/ui/button'
import { Input } from '~/components/ui/input'

export default function Page() {
  return (
    <Form action="/api/search">
      {/* On submission, the input value will be appended to 
          the URL, e.g. /search?query=abc */}
      <Input name="query" />
      <Button type="submit">Search</Button>
    </Form>
  )
}
