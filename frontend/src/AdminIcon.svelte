<script>
  export let value = ''
  export let placeholder = ''

  let fileInput = null

  function onFileChange(e) {
    const file = e.target.files[0]
    if (file.size > 5120) return // 5kb max

    const reader = new FileReader()
    reader.readAsDataURL(file)
    reader.onload = () => {
      value = reader.result
    }
  }
</script>

<button on:click={() => fileInput.click()}>
  <img src={value || placeholder} alt="" />
</button>

<input
  type="file"
  bind:this={fileInput}
  on:change={onFileChange}
  accept="image/*"
/>

<style>
  button {
    width: 26px;
    height: 26px;
    position: relative;
    padding: 0;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  img {
    width: 26px;
    height: 26px;
    object-fit: cover;
    border-radius: 6px;
  }
  input {
    display: none;
  }

  @keyframes rotating {
    0% {
      transform: rotateZ(0deg);
    }
    100% {
      transform: rotateZ(360deg);
    }
  }
</style>
