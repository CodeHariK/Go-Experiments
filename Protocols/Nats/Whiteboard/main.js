import './style.css'
import Alpine from "alpinejs";
import throttle from './throttle';
import { connect, consumerOpts, StringCodec, credsAuthenticator, headers, JSONCodec } from 'nats.ws';

Alpine.data("whiteboard", (subject) => ({
  id: Math.random().toString(36).slice(2, 10),
  color: "black",
  thickness: 5,
  drawing: false,
  last: { x: 0, y: 0 },
  context: null,
  nats: null,
  jc: null,

  async init() {

    console.log(subject)

    this.jc = JSONCodec()

    const creds = await fetch("NGS-Default-CLI.creds")
    if (!creds.ok) {
      addEntry("unable to find NGS-Default-CLI.creds - aborting")
      return;
    }
    const sc = StringCodec()
    const token = await creds.text()
    const auth = credsAuthenticator(sc.encode(token))

    this.nats = await connect({ servers: 'wss://connect.ngs.global', authenticator: auth, debug: true })

    console.log(this.nats)

    const opts = consumerOpts()
    opts.orderedConsumer()
    const sub = await this.nats.jetstream().subscribe(subject, opts)

    for await (const m of sub) {
      const data = this.jc.decode(m.data)
      switch (data.type) {
        case "draw":
          if (data.id !== this.id) {
            this.drawRaw(data)
          }
          break;
        case "clear":
          this.context.clearRect(0, 0, window.innerWidth, window.innerHeight)
        default:
          break;
      }
    }
  },

  sizeCanvas(canvas) {
    this.context = canvas.getContext("2d")
    canvas.width = window.innerWidth
    canvas.height = window.innerHeight
  },

  startDrawing(e) {
    this.drawing = true
    this.last = this.getPoint(e)
  },

  draw(e) {
    throttle(() => {
      const from = this.last
      const to = this.getPoint(e)
      const msg = {
        id: this.id,
        type: "draw",
        from: from,
        to: to,
        thickness: this.thickness,
        color: this.color
      }

      this.drawRaw(msg)
      this.nats.publish(subject, this.jc.encode(msg))

      this.last = to
    }, 30)()
  },

  getPoint(e) {
    if (!e.offsetX || !e.offsetY) {
      const rect = e.target.getBoundingClientRect()
      e.offsetX = (e.touches[0].clientX - window.scrollX - rect.left)
      e.offsetY = (e.touches[0].clientY - window.scrollY - rect.top)
    }
    return { x: e.offsetX, y: e.offsetY }
  },

  clear() {
    const msg = { id: this.id, type: "clear", }
    const h = headers()
    h.set("Nats-Rollup", "sub")
    this.nats.publish(subject, this.jc.encode(msg), { headers: h })
  },

  drawRaw({ from, to, thickness, color }) {
    const c = this.context
    c.beginPath()
    c.lineWidth = thickness
    c.lineCap = "round"
    c.lineJoin = "round"
    c.strokeStyle = color
    c.moveTo(from.x, from.y)
    c.lineTo(to.x, to.y)
    c.stroke()
  },
}))

Alpine.start()
