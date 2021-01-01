const getIconOnType = (type) => {
    switch (type) {
        case "youtube":
            return "subscriptions"
        case "mp3-list":
            return "library_music"
        case "file":
        case "mp3":
            return "queue_music"
    }
    return ""
}
export default {
    state: {
        history: [],
        last_fetch: false,
    },
    mutations: {
        historyfetch(state, data) {
            if (data.pageix === 0) {
                state.history = []
            }
            data.history.forEach(itemsrc => {
                let item = {
                    id: itemsrc.id,
                    icon: getIconOnType(itemsrc.type),
                    playedAt: itemsrc.playedAt,
                    title: itemsrc.title,
                    uri: itemsrc.uri,
                    duration: itemsrc.durationstr,
                }
                state.history.push(item)
            });
            state.last_fetch = (data.history.length === 0)
        }
    }
}