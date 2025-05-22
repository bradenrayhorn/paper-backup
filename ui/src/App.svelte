<script lang="ts">
  import NotFoundPageState from "./lib/NotFoundPageState.svelte";
  import Router from "./lib/router/Router.svelte";
  import type { Route } from "./lib/router/Router.svelte";
  import CreateFileBackupPage from "./pages/CreateFileBackupPage.svelte";
  import CreditsPage from "./pages/CreditsPage.svelte";
  import DecryptFileBackupPage from "./pages/DecryptFileBackupPage.svelte";
  import LandingPage from "./pages/LandingPage.svelte";
  import ShamirSecretCombineFromTextPage from "./pages/ShamirSecretCombineFromTextPage.svelte";
  import ShamirSecretCombinePage from "./pages/ShamirSecretCombinePage.svelte";
  import ShamirSecretCreatePage from "./pages/ShamirSecretCreatePage.svelte";

  const routes: Route[] = [
    { path: "/", component: LandingPage },
    { path: "backup/file", component: CreateFileBackupPage },
    { path: "restore/file", component: DecryptFileBackupPage },
    { path: "shamir/split", component: ShamirSecretCreatePage },
    { path: "shamir/combine", component: ShamirSecretCombinePage },
    { path: "shamir/combine-text", component: ShamirSecretCombineFromTextPage },
    { path: "credits", component: CreditsPage },
    { path: ".*", component: NotFoundPageState },
  ];
</script>

<main>
  <Router {routes} />
</main>
