import Navbar from './components/navbar.js'
import store from './store/index.js'
import routes from './routes.js'


export const app = new Vue({
	el: '#app',
	router: new VueRouter({ routes }),
	components: { Navbar },
	vuetify: new Vuetify(),
	store,
	data() {
		return {
			Buildnr: "",
			links: routes,
			AppTitle: "Omx Control",
			drawer: false,
			connection: null,
		}
	},
	computed: {
		...Vuex.mapState({

		})
	},
	created() {
		// keep in mind that all that is comming from index.html is a string. Boolean or numerics need to be parsed.
		this.Buildnr = window.myapp.buildnr
		let port = location.port;
		let prefix = (window.location.protocol.match(/https/) ? 'wss' : 'ws')
		let socketUrl = prefix + "://" + location.hostname + (port ? ':' + port : '') + "/websocket";
		this.connection = new WebSocket(socketUrl)
		console.log("WS socket created")

		this.connection.onmessage = (event) => {
			console.log(event)
			let dataMsg = JSON.parse(event.data)
			if (dataMsg.type === "status") {
				console.log('Socket msg type: status')
				this.$store.commit('playerstate', dataMsg)
			} else {
				console.warn('Socket message type not recognized ', dataMsg, dataMsg.type)
			}
		}

		this.connection.onopen = (event) => {
			console.log(event)
			console.log("Socket connection success")
		}

		this.connection.onclose = (event) => {
			console.log(event)
			console.log("Socket closed")
			this.connection = null
		}
	},
	methods: {

	},
	template: `
  <v-app class="grey lighten-4">
    <Navbar />
    <v-content class="mx-4 mb-4">
      <router-view></router-view>
    </v-content>
    <v-footer>
      <div class="caption">
        {{ new Date().getFullYear() }} —
        <span>Buildnr: {{Buildnr}}</span>
      </div>
    </v-footer>
  </v-app>`
})

console.log('Main is here!')