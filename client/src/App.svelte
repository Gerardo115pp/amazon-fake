<script>
    import Router from 'svelte-spa-router'
    import { routes } from './routes';
</script>

<style>
    @import url('https://fonts.googleapis.com/css2?family=Rhodium+Libre&display=swap');

    :global(:root) {
        --theme-color: #00a8ff;
        --dimtheme-color: #55c3fa;
        --theme-gradiant: linear-gradient(315deg, rgba(0,249,255,1) 0%, rgba(0,168,255,1) 100%);
        --dark-color: #ffffff;
        --danger:#e72d2d;
        --ready:#25ff3d;
        --modal-shadows: #0000001a;
        --main-font: 'Rhodium Libre';

    }

    :global(body) {
        font-family: 'Rhodium Libre';
        overflow-y: hidden;
        background-color: var(--dark-color);
        color: var(--theme-color);
        margin: 0;
    }

    #app-container {
        width: 100%;
        height: 100vh;
    }
</style>

<div id="app-container">
    <Router {routes}/>
</div>