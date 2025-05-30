// JavaScript for infinite scroll and search on home.html

document.addEventListener('DOMContentLoaded', function () {
    const gamesList = document.getElementById('gamesList');
    const searchForm = document.getElementById('searchForm');
    const searchInput = document.getElementById('searchInput');
    const loading = document.getElementById('loading');
    const noMore = document.getElementById('noMore');
    let page = 1;
    let loadingGames = false;
    let noMoreGames = false;
    let currentQuery = '';

    function createGameCard(game) {
        return /* html */`
        <div class="flex items-start gap-4 p-4 rounded-xl bg-base-200 shadow">
            <div class="w-16 h-16 rounded-lg flex-shrink-0 overflow-hidden">
            ${game.coverImage
                ? /* html */`<img src="${game.coverImage}" alt="${game.title}" class="w-full h-full object-cover">`
                : /* html */`<div class="w-full h-full bg-gradient-to-br from-primary to-secondary"></div>`}
            </div>
            <div>
            <div class="font-bold text-lg">${game.title || 'No Title'}</div>
            <div class="text-sm text-base-content/70">${game.description || 'No description.'}</div>
            </div>
        </div>`;
    }

    async function fetchGames(query, pageNum) {
        loading.classList.remove('hidden');
        try {
            const res = await fetch(`/games/search?title=${encodeURIComponent(query)}&page=${pageNum}`);
            if (res.status === 404) {
                if (pageNum === 1) gamesList.innerHTML = '';
                noMore.classList.remove('hidden');
                loading.classList.add('hidden');
                noMoreGames = true;
                return [];
            }
            if (!res.ok) throw new Error('Failed to fetch games');
            const data = await res.json();
            loading.classList.add('hidden');
            noMore.classList.add('hidden');
            return data.data || data.Data || [];
        } catch (e) {
            loading.classList.add('hidden');
            if (pageNum === 1) gamesList.innerHTML = /* html */`<div class="text-center text-error">Error loading games.</div>`;
            return [];
        }
    }

    async function loadMoreGames() {
        if (loadingGames || noMoreGames || !currentQuery) return;
        loadingGames = true;
        const games = await fetchGames(currentQuery, page);
        if (games.length === 0) {
            noMoreGames = true;
            if (page === 1) noMore.classList.remove('hidden');
        } else {
            gamesList.insertAdjacentHTML('beforeend', games.map(createGameCard).join(''));
            page++;
        }
        loadingGames = false;
    }

    function resetList() {
        gamesList.innerHTML = '';
        page = 1;
        noMoreGames = false;
        noMore.classList.add('hidden');
    }

    searchForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        currentQuery = searchInput.value.trim();
        resetList();
        if (currentQuery) {
            await loadMoreGames();
        }
    });

    window.addEventListener('scroll', async () => {
        if (window.innerHeight + window.scrollY >= document.body.offsetHeight - 200) {
            await loadMoreGames();
        }
    });
});
