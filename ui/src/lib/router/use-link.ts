import type { ActionReturn } from "svelte/action";

type Parameter = undefined | (() => void);

export function link(
  node: HTMLElement,
  callback?: Parameter,
): ActionReturn<Parameter> {
  function onClick(e: MouseEvent) {
    if (callback) {
      callback();
    }
    e.preventDefault();
    window.history.pushState(undefined, "", node.getAttribute("href"));
  }

  return {
    update: () => {
      node.removeEventListener("click", onClick);
      node.addEventListener("click", onClick);
    },
    destroy: () => {
      node.removeEventListener("click", onClick);
    },
  };
}
