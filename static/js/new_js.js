const pdfEditor = {
    CANVAS_ELEMENTS: [],
    CLICKED_ON_CANVAS: false,
    CLICKED_ON_BEFORE: false,
    HOST: "",
    COLOR_CHAGED: "",
    init: function (host) {
        this.HOST = host;
        this.initDropZone();
        $("#addbutton").on("click", pdfEditor.addText);
        $("#save").on("click", pdfEditor.savePdf);
        $("#color-control").spectrum({
            color: "#000",
            preferredFormat: "hex",
            showInput: true,
            showPalette: true,
            palette: [["red", "rgba(0, 255, 0, .5)", "rgb(0, 0, 255)"]],
            change: function (color) {
                console.log("change");
                pdfEditor.COLOR_CHAGED = color.toHexString();

                console.log(pdfEditor.COLOR_CHAGED);

                let canvasData = pdfEditor.getSelectedCanvas();
                if (canvasData !== null) {
                    let CANVAS = canvasData.canvas;

                    if (CANVAS.getActiveObject() !== null && CANVAS.getActiveObject() !== undefined) {
                        CANVAS.getActiveObject().set('fill', pdfEditor.COLOR_CHAGED);
                        CANVAS.renderAll();
                    }
                }
            }
        });

        $("#page").on("change", function (e) {
            let canvasData = pdfEditor.getSelectedCanvas();
            let CANVAS = canvasData.canvas;

            CANVAS.on("object:selected", function (e) {
                pdfEditor.canvasListener(CANVAS, CANVAS.getActiveObject());
            });
        });

        $(".pdf_data_val").on("click", function (e) {
            //dublicate
            let text = new fabric.Text(e.target.text, {
                fontSize: 24,
                fontWeight: parseInt($('#fontWeight-control').val(), 10),
                charSpacing: parseInt($('#charSpacing-control').val(), 10),
            });

            let color = $('#color-control').val();

            if (color === "") {
                color = "#000";
            } else {
                if (color[0] !== "#") {
                    color = "#" + color;
                }
            }

            text.setColor(color);

            let canvasData = pdfEditor.getSelectedCanvas();
            let CANVAS = canvasData.canvas;

            CANVAS.add(text);

            CANVAS.on("selection:created", function (e) {
                $('html').keyup(function (e) {
                    if (e.keyCode === 46) {
                        pdfEditor.deleteElement();
                    }
                });
            });

            $('html').keyup(function (e) {
                if (e.keyCode === 46) {
                    pdfEditor.deleteElement();
                }
            });
        });
    },
    initDropZone: function () {
        Dropzone.options.pdfAwesome = {
            maxFiles: 1,
            url: `${pdfEditor.HOST}/api/v1/pdf/upload`,
            addRemoveLinks: false,
            acceptedFiles: ".pdf",
            dictDefaultMessage: "Drop or select pdf file",
            accept: function (file, done) {
                console.log("uploaded");
                done();
            },
            init: function () {
                this.on("maxfilesexceeded", function (file) {
                    console.log(file);
                });

                this.on('addedfile', function (file) {
                    if (this.files.length > 1) {
                        this.removeFile(this.files[0]);
                    }
                });

                this.on("sending", function (file) {
                    console.log('upload started', file);
                    $('.meter').show();
                });

                // File upload Progress
                this.on("totaluploadprogress", function (progress) {
                    console.log("progress ", progress);
                    $('.roller').width(progress + '%');
                });

                this.on("queuecomplete", function (progress) {
                    $('.meter').delay(999).slideUp(999);
                });

                this.on("success", function (file, responseText) {
                    if (responseText.images.length > 0) {
                        $(".dropzone").animate({"minHeight": 150});
                        $("#secondStep").show();
                        let $canvasArea = $('#canvasArea');
                        $canvasArea.empty();
                        let i = 1;
                        let $pageSelect = $("#page");
                        for (let image of responseText.images) {
                            let id = `canvas_${i}`;
                            let canvas = pdfEditor.createCanvasElement(id);
                            $canvasArea.append(canvas);
                            $pageSelect.append(`<option value="${i}">Page ${i}</option>`);

                            pdfEditor.initPdfEditor(canvas, image, i);
                            i++;
                        }
                        $("#pdfAwesome").delay(300).hide();
                        $pageSelect.val(1).trigger('change');


                        let canvasData = pdfEditor.getSelectedCanvas();
                        let CANVAS = canvasData.canvas;

                        CANVAS.on("object:selected", function (e) {
                            if (CANVAS.getActiveObject() !== null) {
                                if (CANVAS.getActiveObject().text) {
                                    pdfEditor.canvasListener(CANVAS, CANVAS.getActiveObject());
                                }
                            }
                        });
                    }
                });
            }
        };
    },
    createCanvasElement: function (id) {
        this.canvas = document.createElement('canvas');
        this.canvas.id = id;
        this.canvas.width = 828;
        this.canvas.height = 1100;
        return this.canvas;
    },
    initPdfEditor: function (canvas, src, pageNumber) {
        let CANVAS = new fabric.Canvas(canvas);
        fabric.Object.prototype.transparentCorners = false;

        CANVAS.setBackgroundImage(
            src,
            CANVAS.renderAll.bind(CANVAS),
            {
                crossOrigin: 'anonymous'
            }
        );

        pdfEditor.CANVAS_ELEMENTS.push({"page": pageNumber, "canvas": CANVAS, "src": src});

        CANVAS.on('mouse:down', function (options) {
            if (options.target === null) {
                if (!pdfEditor.CLICKED_ON_CANVAS) {
                    pdfEditor.addTextByClick(CANVAS, options.e);
                    pdfEditor.CLICKED_ON_BEFORE = true;
                    pdfEditor.CLICKED_ON_CANVAS = true;
                } else {
                    if (pdfEditor.CLICKED_ON_BEFORE) {
                        pdfEditor.CLICKED_ON_CANVAS = false;
                        pdfEditor.CLICKED_ON_BEFORE = false;
                    }
                }
            }
        });

        return CANVAS
    },
    getSelectedCanvas: function () {
        const page = $('#page').val();
        let resultCanvas = null;

        if (pdfEditor.CANVAS_ELEMENTS.length > 0) {
            for (let canvas of pdfEditor.CANVAS_ELEMENTS) {
                console.log(canvas.page);
                console.log(canvas.page == page);
                if (canvas.page == page) {
                    resultCanvas = canvas;
                    break;
                }
            }
        }

        return resultCanvas;
    },
    addText: function (e) {
        let $newText = $('#newtext');
        let text = new fabric.Text($newText.val(), {
            fontSize: 24,
            fontWeight: parseInt($('#fontWeight-control').val(), 10),
            charSpacing: parseInt($('#charSpacing-control').val(), 10),
        });

        let color = $('#color-control').val();

        if (color === "") {
            color = "#000";
        } else {
            if (color[0] !== "#") {
                color = "#" + color;
            }
        }

        text.setColor(color);

        let canvasData = pdfEditor.getSelectedCanvas();
        let CANVAS = canvasData.canvas;

        CANVAS.add(text);

        CANVAS.on("selection:created", function (e) {
            $('html').keyup(function (e) {
                if (e.keyCode === 46) {
                    pdfEditor.deleteElement();
                }
            });
        });

        $('html').keyup(function (e) {
            if (e.keyCode === 46) {
                pdfEditor.deleteElement();
            }
        });
    },
    addTextByClick: function (canvas, e) {
        const text = new fabric.IText('...', {
            fontSize: 24,
            fontWeight: parseInt($('#fontWeight-control').val(), 10),
            charSpacing: parseInt($('#charSpacing-control').val(), 10),
            top: e.offsetY,
            cursorDuration: 500,
            left: e.offsetX
        });

        let color = $('#color-control').val();

        if (color === "") {
            color = "#000";
        } else {
            if (color[0] !== "#") {
                color = "#" + color;
            }
        }

        text.setColor(color);

        canvas.add(text);
    },
    canvasListener: function (CANVAS, elem) {
        console.log(CANVAS.getActiveObject());
        if (CANVAS.getActiveObject() === null) {
            return;
        }

        const $angleControl = $('#angle-control');
        $angleControl.on("input", function () {
            CANVAS.getActiveObject().set('angle', parseInt(this.value, 10)).setCoords();
            CANVAS.renderAll();
        });

        const $scaleControl = $('#scale-control');
        $scaleControl.on("input", function () {
            CANVAS.getActiveObject().scale(parseFloat(this.value)).setCoords();
            CANVAS.renderAll();
        });

        const $topControl = $('#top-control');
        $topControl.on("input", function () {
            CANVAS.getActiveObject().set('top', parseInt(this.value, 10)).setCoords();
            CANVAS.renderAll();
        });

        const $leftControl = $('#left-control');
        $leftControl.on("input", function () {
            CANVAS.getActiveObject().set('left', parseInt(this.value, 10)).setCoords();
            CANVAS.renderAll();
        });

        const $skewXControl = $('#skewX-control');
        $skewXControl.on("input", function () {
            CANVAS.getActiveObject().set('skewX', parseInt(this.value, 10)).setCoords();
            CANVAS.renderAll();
        });

        const $skewYControl = $('#skewY-control');
        $skewYControl.on("input", function () {
            CANVAS.getActiveObject().set('skewY', parseInt(this.value, 10)).setCoords();
            CANVAS.renderAll();
        });

        const $charSpacingControl = $('#charSpacing-control');
        $charSpacingControl.on("input", function () {
            CANVAS.getActiveObject().set('charSpacing', parseInt(this.value, 10)).setCoords();
            CANVAS.renderAll();
        });

        const $colorControl = $('#color-control');
        $colorControl.on("change", function () {
            if (CANVAS.getActiveObject() !== null) {
                CANVAS.getActiveObject().set('fill', this.value);
                CANVAS.renderAll();
            }
        });

        const $fontWeightControl = $('#fontWeight-control');
        $fontWeightControl.on("change", function () {
            CANVAS.getActiveObject().set('fontWeight', parseInt(this.value, 10));
            CANVAS.renderAll();
        });

        function updateControls() {
            $scaleControl.val(elem.scaleX);
            $angleControl.val(elem.angle);
            $leftControl.val(elem.left);
            $topControl.val(elem.top);
            $skewXControl.val(elem.skewX);
            $skewYControl.val(elem.skewY);
            $charSpacingControl.val(elem.charSpacing);
            $colorControl.val(elem.fill.substring(1));
            $fontWeightControl.val(elem.fontWeight);
        }

        CANVAS.on({
            'object:moving': updateControls,
            'object:scaling': updateControls,
            'object:resizing': updateControls,
            'object:rotating': updateControls,
            'object:skewing': updateControls,
            'object:modified': updateControls
        });
    },
    deleteElement: function () {
        let canvasData = pdfEditor.getSelectedCanvas();
        let CANVAS = canvasData.canvas;

        const activeObject = CANVAS.getActiveObject(),
            activeGroup = CANVAS.getActiveGroup();

        if (activeObject) {
            if (confirm('Are you sure?')) {
                CANVAS.remove(activeObject);
            }
        } else if (activeGroup) {
            if (confirm('Are you sure?')) {
                const objectsInGroup = activeGroup.getObjects();
                CANVAS.discardActiveGroup();
                objectsInGroup.forEach(function (object) {
                    CANVAS.remove(object);
                });
            }
        }
    },
    savePdf: function () {
        let results = [];
        let i = 1;

        for (let page of pdfEditor.CANVAS_ELEMENTS) {
            results.push({"page": i, "b64": page.canvas.toDataURL()});
            i++;
        }

        $.ajax({
            "method": "POST",
            "url": `${pdfEditor.HOST}/api/v1/pdf/save`,
            "data": JSON.stringify(results)
        }).done(function (data) {
            console.log(data);
            window.open(data.pdf, '_blank');
        }).fail(function () {

        });
    }
};

pdfEditor.init("https://localhost:8443");