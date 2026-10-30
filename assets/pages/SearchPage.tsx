import { useLocation } from "preact-iso";
import { withViewTransition } from "../functions/dom.ts";
import { useSearchResults } from "../hooks/useSearchResults.tsx";
import { Timer } from "../components/Timer.tsx";
import { SearchForm } from "../components/search/SearchForm.tsx";
import { SearchItem } from "../components/search/SearchItem.tsx";
import { ImageItem } from "../components/search/ImageItem.tsx";
import { SearchToggle } from "../components/search/SearchToggle.tsx";
import { SearchWallpaperButton } from "../components/search/SearchWallpaperButton.tsx";
import { useArrowNavigation } from "../hooks/useArrowNavigation.ts";

export function SearchPage() {
  const location = useLocation();
  const { columns, imageColumns, query, searchMode, isFetching, component } = useSearchResults(location);
  useArrowNavigation();

  const onSearch = (q: string) => {
    const url = new URL(window.location.pathname, window.location.origin);
    url.searchParams.set("q", q);
    // Preserve search mode
    if (searchMode === "images") {
      url.searchParams.set("mode", "images");
    }
    withViewTransition(() => location.route(url.toString()));
  };

  const onModeChange = (mode: "web" | "images") => {
    const url = new URL(window.location.pathname, window.location.origin);
    if (query) {
      url.searchParams.set("q", query);
    }
    if (mode === "images") {
      url.searchParams.set("mode", "images");
    } else {
      url.searchParams.delete("mode");
    }
    withViewTransition(() => location.route(url.toString()));
  };

  const showWallpaperSwitcher = !query;

  return (
    <>
      <header class="search-top">
        <SearchForm
          onSearch={onSearch}
          defaultValue={query}
          isLoading={isFetching}
        />
        {query && (
          <SearchToggle 
            onModeChange={onModeChange} 
            defaultMode={searchMode as "web" | "images"}
          />
        )}
        {component.value}
        <Timer />
      </header>
      {showWallpaperSwitcher && <SearchWallpaperButton />}
      <main class={`search-main ${searchMode === "images" ? "images-mode" : ""}`}>
        {searchMode === "images" ? (
          // Image results
          imageColumns.value.map((column, index) => (
            <div class="search-column" key={index}>
              {column.map((result, k) => (
                <ImageItem result={result} key={k} />
              ))}
            </div>
          ))
        ) : (
          // Web results
          columns.value.map((column, index) => (
            <div class="search-column" key={index}>
              {column.map((result, k) => (
                <SearchItem result={result} key={k} />
              ))}
            </div>
          ))
        )}
      </main>
    </>
  );
}
