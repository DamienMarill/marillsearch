import { useSignal } from "@preact/signals";
import { useEffect } from "preact/hooks";
import clsx from "clsx";

type SearchMode = "web" | "images";

type Props = {
  onModeChange: (mode: SearchMode) => void;
  defaultMode?: SearchMode;
};

export function SearchToggle({ onModeChange, defaultMode = "web" }: Props) {
  const currentMode = useSignal<SearchMode>(defaultMode);

  // Sync currentMode with defaultMode changes
  useEffect(() => {
    currentMode.value = defaultMode;
  }, [defaultMode]);

  const handleToggle = (mode: SearchMode) => {
    currentMode.value = mode;
    onModeChange(mode);
  };

  return (
    <div class="search-toggle">
      <button
        type="button"
        class={clsx("search-toggle-btn", currentMode.value === "web" && "active")}
        onClick={() => handleToggle("web")}
      >
        Web
      </button>
      <button
        type="button"
        class={clsx("search-toggle-btn", currentMode.value === "images" && "active")}
        onClick={() => handleToggle("images")}
      >
        Images
      </button>
    </div>
  );
}