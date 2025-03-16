<script lang="ts">
  import { onDestroy, onMount } from "svelte";

  type Camera = { id: string; label: string };

  const {
    onError,
    onFrame,
  }: {
    onError: (error: Error) => void;
    onFrame: (frame: Uint8ClampedArray, width: number, height: number) => void;
  } = $props();

  let stream: MediaStream | null = null;

  let displayCanvas = $state<HTMLCanvasElement>();
  let scanCanvas = $state<HTMLCanvasElement>();
  let video = $state<HTMLVideoElement>();
  let cameras = $state<Camera[]>([]);

  async function initDevices() {
    const device = await navigator.mediaDevices.getUserMedia({ video: true });
    device.getTracks().forEach((track) => {
      track.stop();
    });

    const devices = await navigator.mediaDevices.enumerateDevices();

    cameras = devices
      .filter((device) => device.kind === "videoinput")
      .map((device) => ({ id: device.deviceId, label: device.label }));
  }

  function drawVideoFrame(
    ctxDisplay: CanvasRenderingContext2D,
    ctxScan: CanvasRenderingContext2D,
  ) {
    if (!video || !displayCanvas || !scanCanvas || video.readyState !== 4)
      return;

    // draw frame to display and scan canvas
    ctxDisplay.drawImage(
      video,
      0,
      0,
      displayCanvas.width,
      displayCanvas.height,
    );
    ctxScan.drawImage(video, 0, 0, scanCanvas.width, scanCanvas.height);

    // read data from scan canvas
    const imageData = scanCanvas
      .getContext("2d")
      ?.getImageData(0, 0, scanCanvas.width, scanCanvas.height);
    if (imageData) {
      onFrame(imageData.data, scanCanvas.width, scanCanvas.height);
    }

    requestAnimationFrame(() => drawVideoFrame(ctxDisplay, ctxScan));
  }

  function resizeVideo() {
    if (!video || !displayCanvas || !scanCanvas) return;

    // scan canvas is max resolution
    scanCanvas.width = video.videoWidth;
    scanCanvas.height = video.videoHeight;

    // calculate display canvas disze
    const maxWidth =
      displayCanvas.parentElement?.clientWidth ?? window.innerWidth;
    const maxHeight = window.innerHeight - 40;

    const videoAspectRatio = video.videoWidth / video.videoHeight;

    let canvasWidth = maxWidth;
    let canvasHeight = canvasWidth / videoAspectRatio;

    if (canvasHeight > maxHeight) {
      canvasHeight = maxHeight;
      canvasWidth = canvasHeight * videoAspectRatio;
    }

    displayCanvas.width = canvasWidth;
    displayCanvas.height = canvasHeight;

    // start drawing
    const ctxDisplay = displayCanvas.getContext("2d");
    const ctxScan = scanCanvas.getContext("2d", { willReadyFrequently: true });
    if (ctxDisplay && ctxScan)
      drawVideoFrame(
        ctxDisplay as CanvasRenderingContext2D,
        ctxScan as CanvasRenderingContext2D,
      );
  }

  async function onSelectCamera(id: string | null) {
    stream?.getTracks().forEach((track) => track.stop());
    if (displayCanvas) {
      displayCanvas.height = 0;
      displayCanvas.width = 0;
    }
    if (id === null) {
      return;
    }

    stream = await navigator.mediaDevices.getUserMedia({
      video: { deviceId: id, width: { ideal: 4096 }, height: { ideal: 2160 } },
    });
    if (!video) return;

    video.srcObject = stream;
    video.play();
  }

  onMount(() => {
    initDevices().catch((error) => {
      onError(error);
    });
  });

  onDestroy(() => {
    stream?.getTracks().forEach((track) => track.stop());
  });
</script>

<div>
  <select
    onchange={(e) => {
      onSelectCamera(
        e.currentTarget.value === "none" ? null : e.currentTarget.value,
      ).catch((error) => {
        onError(error);
      });
    }}
  >
    <option value="none">None</option>
    {#each cameras as camera (camera.id)}
      <option value={camera.id}>{camera.label}</option>
    {/each}
  </select>

  <canvas bind:this={displayCanvas}></canvas>
  <canvas bind:this={scanCanvas} style="display: none;"></canvas>
</div>

<video
  bind:this={video}
  playsinline={true}
  style="display: none;"
  aria-hidden="true"
  onloadedmetadata={resizeVideo}
></video>

<style>
  div {
    display: block;
  }

  select {
    width: 100%;
  }

  canvas {
    margin-top: 0.5rem;
  }
</style>
