<script>
    import Button from '../../../components/Button.svelte';
    import { onMount } from "svelte";
import { server_name } from '../../../server_info';
import { pop } from 'svelte-spa-router';

    export let session_key = "";
    export let products_data = [];

    let total = 0;
    $: if (products_data.length > 0) {
            total = 0;
            products_data.forEach(p => {
                total += p.unit_price * p.count;
            });
        }

    const completePurchase = () => {
        const headers = new Headers();
        headers.set("X-sk", session_key);
        
        const request = new Request(`http://${server_name}/cart`, {method: 'PATCH', headers: headers});
        fetch(request)
            .then(promise => {
                if (promise.ok) {
                    alert("Purchase completed!");
                } else {
                    alert("Purchase failed");
                }
                pop();
            });
    }

</script>

<style>

    @keyframes SpawnCheckout {
        0% {
            transform: translate(0, 100%);
            opacity: 20%;
        }
        100% {
            transform: translate(0, 0);
            opacity: 100%;
        }
    }

    #check-out-component {
        color: white;
        padding: 0 2vw;
        animation-name: SpawnCheckout;
        animation-duration: 1s;
        animation-iteration-count: 1;
        
    }

    #checkout-title {
        font-size: 1.9rem;
        text-transform: uppercase;
        text-align: center;
        border-bottom: 2px solid white;
    }

    #checkout-details {
        height: 30vh;
        border-bottom: 2px solid white;
        padding: 0 1vw;
    }

    .checkout-detail-container {
        display: flex;
        justify-content: space-between;
    }

    #checkout-conclusions {
        padding: 0 1vw;
        border-bottom: 2px solid white;
    }
    
    .cc-conclusion {
        display: block;
    }

    .cc-total {
        font-size: 1.2em;
        font-weight: bold;
    }

    #checkout-controls {
        display: flex;
        height: 16vh;
        justify-content: center;
        align-items: center;
    }
</style>

<div id="check-out-component">
    <div id="checkout-title">
        Checkout
    </div>
    <div id="checkout-details">
        {#each products_data as pd}
            <div class="checkout-detail-container">
                <div class="cd-name">{pd.product_name}</div>
                <div class="cd-amount">
                    <span>{pd.count} x </span>
                    <span>${Number(pd.unit_price).toLocaleString('en')} MXN</span>
                </div>
            </div>
        {/each}
    </div>
    <div id="checkout-conclusions">
        <span class="cc-conclusion cc-before-iva">Before iva: ${total.toLocaleString('en')} MXN</span>
        <span class="cc-conclusion cc-iva">Iva: ${(total*0.16).toLocaleString('en')} MXN</span>
        <span class="cc-conclusion cc-total">Total: ${(total + (total*0.16)).toLocaleString('en')} MXN</span>
    </div>
    <div id="checkout-controls">
        <Button 
            width="40%"
            button_color="var(--ready)"
            label="COMPLETE PURCHASE"
            onClick={completePurchase}
        />
    </div>
</div>