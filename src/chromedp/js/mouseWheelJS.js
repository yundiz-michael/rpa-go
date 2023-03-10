function mouseWheel(offsetY) {
    const evt = new Event('wheel', {bubbles: true, cancelable: true});
    evt.deltaY += offsetY;
    this.dispatchEvent(evt);
    return "true";
}
