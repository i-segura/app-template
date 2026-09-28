import { useEffect, useState } from "react";

interface HelloResponse {
  message: string;
}

export default function App() {
  const [message, setMessage] = useState<string>("loading…");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch("/api/v1/hello")
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json() as Promise<HelloResponse>;
      })
      .then((data) => setMessage(data.message))
      .catch((err: unknown) =>
        setError(err instanceof Error ? err.message : String(err)),
      );
  }, []);

  return (
    <main>
      <h1>{{ cookiecutter.project_name }}</h1>
      <p>{{ cookiecutter.description }}</p>
      <section aria-live="polite">
        {error ? <p role="alert">API error: {error}</p> : <p>{message}</p>}
      </section>
    </main>
  );
}
