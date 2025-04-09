<script lang="ts">
  import { link } from "../lib/router/use-link";
</script>

Pick a tool.

<ul>
  <li><a use:link href="/backup/file">Create File Backup</a></li>
  <li><a use:link href="/restore/file">Restore File Backup</a></li>
  <li><a use:link href="/shamir/split">Shamir Secret Split</a></li>
  <li><a use:link href="/shamir/combine">Shamir Secret Combine</a></li>
  <li>
    <a use:link href="/shamir/combine-text">Shamir Secret Combine From Text</a>
  </li>
</ul>
