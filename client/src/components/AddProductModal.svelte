<script>
    export let close = () => {};

    const handleBackgroundClick = e => {
        if (e.target === e.currentTarget) {
            close();
        }
    }
</script>

<style>
    #add-product-background {
        display: flex;
        position: fixed;
        width: 100vw;
        height: 100vh;
        background-color: rgba(0, 0, 0, 0.2);
        justify-content: center;
    }

    #ap-modal {
        width: 80vw;
        background-color: white;
        border-left: 1px solid var(--theme-color);
        border-right: 1px solid var(--theme-color);
    }
</style>

<div on:click={handleBackgroundClick} id="add-product-background">
    <div id="ap-modal">A</div>
</div>