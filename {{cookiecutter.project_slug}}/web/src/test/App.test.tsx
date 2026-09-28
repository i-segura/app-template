import { render, screen, waitFor } from "@testing-library/react";
import App from "../App";

// Vitest example test: stub fetch, render App, assert the API message shows up.
describe("App", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("renders the greeting returned by the API", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => ({ message: "hello, test" }),
      }),
    );

    render(<App />);
    await waitFor(() =>
      expect(screen.getByText("hello, test")).toBeInTheDocument(),
    );
  });

  it("shows an error when the API fails", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({ ok: false, status: 500 }),
    );

    render(<App />);
    await waitFor(() =>
      expect(screen.getByRole("alert")).toHaveTextContent("HTTP 500"),
    );
  });
});
