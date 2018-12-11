import { Vue } from "vue"

new Vue({
	el: "#app",
	template: "<app></app>"
})

Vue.component("app", {
	render(h) {
		return <div>h</div>
	}
})
