import Login from './pages/Login.svelte';
import SignUp from './pages/SignUp.svelte';
import Store from './pages/StoreMainView.svelte';
import UserProfile from './pages/UserProfile.svelte';


const routes = {
    "/": Login,
    "/sign-up": SignUp,
    "/store/:sk": Store,
    "/profile/:sk": UserProfile 
}

export { routes }