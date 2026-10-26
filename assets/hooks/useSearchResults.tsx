import { useSignal } from "@preact/signals";
import { type ComponentChild } from "preact";
import { type LocationHook } from "preact-iso";
import { useEffect } from "preact/hooks";
import { withViewTransition } from "../functions/dom.ts";
import { jsonFetch } from "../functions/http.ts";
import { matchMiddlewares } from "../functions/middleware/middlewares.tsx";

type SearchResultLink = {
  title: string;
  url: string;
};

export type SearchResult = {
  url: string;
  title: string;
  desc: string;
  domain: string;
  author?: string;
  related?: SearchResultLink[];
  siteName?: string;
};

export type ImageResult = {
  url: string;
  title: string;
  thumbnail: string;
  source: string;
  domain: string;
  width?: number;
  height?: number;
  size?: string;
};

const baseTitle = "GrafiSearch";

export function useSearchResults(location: LocationHook) {
  const columns = useSignal<SearchResult[][]>([]);
  const imageColumns = useSignal<ImageResult[][]>([]);
  const query = location.query.q;
  const searchMode = location.query.mode === "images" ? "images" : "web";
  const isFetching = useSignal(false);
  const component = useSignal<ComponentChild>(null);

  // Update title according to query
  useEffect(() => {
    document.title = query ? `${query} - ${baseTitle}` : baseTitle;
  }, [query]);

  const pushColumn = (column: SearchResult[]) => {
    withViewTransition(() => {
      columns.value = [...columns.value, column];
    });
  };

  const pushImageColumn = (column: ImageResult[]) => {
    withViewTransition(() => {
      imageColumns.value = [...imageColumns.value, column];
    });
  };

  useEffect(() => {
    if (!query) {
      withViewTransition(() => {
        columns.value = [];
        imageColumns.value = [];
        document.body.classList.remove("has-results");
      });
      return;
    }
    document.body.classList.add("has-results");

    // Check if a smart function or bang intercepts the search
    component.value = matchMiddlewares(query);
    if (component.value) {
      return;
    }

    // Clear previous results
    columns.value = [];
    imageColumns.value = [];
    const abortController = new AbortController();
    const signal = abortController.signal;
    isFetching.value = true;

    if (searchMode === "images") {
      // Image search - Pinterest only
      jsonFetch<ImageResult[]>("/api/images/pinterest", {
        query: { q: query },
        signal,
      })
        .then(pushImageColumn)
        .then(() => {
          isFetching.value = false;
        })
        .catch(() => {
          isFetching.value = false;
        });
    } else {
      // Web search
      Promise.allSettled([
        jsonFetch<SearchResult[]>("/api/ddg", {
          query: { q: query },
          signal,
        }).then(pushColumn),
        jsonFetch<SearchResult[]>("/api/google", {
          query: { q: query },
          signal,
        }).then(pushColumn),
      ]).then(() => {
        isFetching.value = false;
      });
    }

    return () => abortController.abort();
  }, [query, searchMode]);

  return { columns, imageColumns, query, searchMode, isFetching, component };
}
