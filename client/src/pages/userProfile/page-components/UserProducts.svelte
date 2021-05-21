<script>
    export let products = [];
</script>

<style>
    #user-products-container {
        background: var(--theme-color);
        width: 20vw;
        height: 60vh;
        color: white;
        padding: 1vh 1vw;
        border-radius: 15px;
    }

    .user-product-container {
        display: flex;
        height: 5vh;
        justify-content: space-around;
        border-bottom: 1px solid white;
    }

    .user-product-container * {
        user-select: none;
    }
</style>

<div id="user-products-container">
    {#each products as p}
        <div class="user-product-container">
            <span class="product-name">{p.name}</span>
            <span class="product-price">{p.price} MXN</span>
            <span class="product-sales">{p.solds}/{p.stock}</span>
        </div>
    {/each}
</div>