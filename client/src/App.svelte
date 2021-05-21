<script>
    import Router from 'svelte-spa-router'
    import { routes } from './routes';
</script>

<style>
    @import url('https://fonts.googleapis.com/css2?family=Rhodium+Libre&display=swap');

    :global(:root) {
        --theme-color: #78ebff;
        --dimtheme-color: #acf5fc;
        --theme-gradiant: linear-gradient(115.58deg, #78ebff 7.24%, #00c3ff 93.27%);;
        --dark-color: #ffffff;
        --danger:#e72d2d;
        --ready:#acfd50;
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