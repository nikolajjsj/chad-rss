import Link from "next/link";

export default function LandingPage() {
  return (
    <main className="min-h-screen">
      <div className="bg-gray-900">
        <div className="container mx-auto px-6 py-20">
          <h1 className="text-4xl font-medium text-white">ChadRSS</h1>
          <p className="mt-2 text-white">A chad RSS reader</p>
          <Link href="/api/auth/signin" className="text-white">
            Get started
          </Link>
        </div>
      </div>
      <div className="container mx-auto px-6 py-20">
        <h2 className="text-2xl font-medium">Features</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-6">
          <div className="flex items-center">
            <svg
              xmlns="http://www.w3.org/
              2000/svg"
              className="h-6 w-6 mr-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M5 13l4 4L19 7"
              />
            </svg>
            <span>Read your favorite feeds</span>
          </div>
          <div className="flex items-center">
            <svg
              xmlns="http://www.w3.org/
              2000/svg"
              className="h-6 w-6 mr-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M5 13l4 4L19 7"
              />
            </svg>
            <span>Save articles for later</span>
          </div>
          <div className="flex items-center">
            <svg
              xmlns="http://www.w3.org/
              2000/svg"
              className="h-6 w-6 mr-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M5 13l4 4L19 7"
              />
            </svg>
            <span>Share articles with friends</span>
          </div>
          <div className="flex items-center">
            <svg
              xmlns="http://www.w3.org/
              2000/svg"
              className="h-6 w-6 mr-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M5 13l4 4L19 7"
              />
            </svg>
            <span>Dark mode</span>
          </div>
        </div>
      </div>
    </main>
  );
}
