import type { ActionReturn } from "svelte/action";

type Parameter = undefined | (() => void);

function addPrefixToHref(node: HTMLAnchorElement) {
  const basePath = window.location.pathname.replace(/\/+$/, "");

  // remove protocol and host from href. these are automatically added by the browser.
  let href = node.href.replace(/^[a-z0-9+.-]+:\/\//i, "");
  if (href.startsWith(window.location.host)) {
    href = href.substring(window.location.host.length);
  }

  if (!href.startsWith(`${basePath}/#/`)) {
    href = href.replace(/^\/+/, "");
    node.href = `${basePath}/#/${href}`;
  }
}

export function link(
  node: HTMLAnchorElement,
  callback?: Parameter,
): ActionReturn<Parameter> {
  function onClick(e: MouseEvent) {
    if (callback) {
      callback();
    }
    e.preventDefault();
    window.history.pushState(undefined, "", node.getAttribute("href"));
  }

  node.addEventListener("click", onClick);
  addPrefixToHref(node);

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
