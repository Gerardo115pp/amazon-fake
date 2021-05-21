<script>
    import Button from '../../../components/Button.svelte';
    import goback_svg from '../../../resources/return.svg';
    import add_svg from '../../../resources/add.svg';
    import { pop, push } from 'svelte-spa-router';
import { server_name } from '../../../server_info';

    export let session_key = "";
    export let closeProductModal = () => {};
    export let user_data = {};

    const logout = () => {
        const headers = new Headers();
        headers.set("X-sk", session_key);

        const request = new Request(`http://${server_name}/logout`, {method: 'PATCH', headers: headers});
        fetch(request).then(promise => {
            if(promise.ok) {
                push("/");
            }
        });
    }
</script>

<style>
    #profile-controls {
        display: flex;
        width: 97vw;
        padding: 0.5% 1%;
        justify-content: space-between;
        align-items: center;
    }

    :global(#profile-controls svg) {
        cursor: pointer;
        width: 3vw;
        fill: var(--theme-color);
    }

    #user-name-label {
        font-size: 1.3rem;
    }

    #add-product {
        margin-left: 75%;
    }

    :global(#add-product svg) {
        width: 2vw;
    }

    #logout-btn-container {
        margin-left: 2%;
    }
</style>

<header id="profile-controls">
    <div on:click={() => pop()} class="goback-container">
        {@html goback_svg}
    </div>
    <div on:click={closeProductModal} id="add-product">
        {@html add_svg}
    </div>
    <div id="user-name-label">
        <span>{user_data.username}</span>
    </div>
    <div id="logout-btn-container">
        <Button onClick={logout} width="5vw" label="logout" />
    </div>
</header>