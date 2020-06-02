const {
    Component,
} = window.Torus;

async function getBusySlots(startTime, endTime) {
    const data = await fetch('/data', {
        method: 'POST',
        body: JSON.stringify({
            timeMin: x,
            timeMax: x,
        })
    }).then(resp => resp.json());
    const calendars = Object.values(data.calendars);
    // TODO
}

class App extends Component {

    compose() {
        return jdom`<div class="app">
            hi
        </div>`;
    }
}

const app = new App();
document.getElementById('root').appendChild(app.node);
