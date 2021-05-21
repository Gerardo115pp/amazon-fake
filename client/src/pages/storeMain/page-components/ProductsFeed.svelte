<script>
    export let products = [];
</script>

<style>
    #products-feed {
        width: 50vw;
        height: 50vh;
    }

    #feed-title {
        border-bottom: 3px solid var(--theme-color);
    }

    .product-container {
        display: flex;
        background: var(--theme-color);
        color: white;
        padding: 1vh 1vw;
        margin: 1vh 0;
        justify-content: space-around;
    }
</style>

<div id="products-feed">
    <div id="feed-title">
        <h2>Store feed: {products.length} products</h2>
    </div>
    <div id="feed-container">
        {#each products as p}
            <div class="product-container">
                <span class="product-name">{p.name}</span>
                <span class="product-stock">{p.solds}/{p.stock}</span>
            </div>
        {/each}
    </div>
</div>