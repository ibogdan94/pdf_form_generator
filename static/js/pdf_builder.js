const pdfEditor = {
    CANVAS_ELEMENTS: [],
    CLICKED_ON_CANVAS: false,
    CLICKED_ON_BEFORE: false,
    init: function () {
        this.initDropZone();
        $("#addbutton").on("click", pdfEditor.addText);
        $("#save").on("click", pdfEditor.savePdf);
        $("#delete").on("click", pdfEditor.deleteElement);
        $("#color-control").spectrum({
            color: "#000",
            preferredFormat: "hex",
            showInput: true,
            showPalette: true,
            palette: [["red", "rgba(0, 255, 0, .5)", "rgb(0, 0, 255)"]]
        });
    },
    initDropZone: function () {
        Dropzone.options.myAwesomeDropzone = {
            maxFiles: 1,
            url: 'api/v1/pdf/upload',
            addRemoveLinks: false,
            acceptedFiles: ".pdf",
            accept: function(file, done) {
                console.log("uploaded");
                done();
            },
            init: function() {
                this.on("maxfilesexceeded", function(file){
                    console.log(file);
                });

                this.on('addedfile', function(file) {
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

                this.on("success", function(file, responseText) {
                    if (responseText.images.length > 0) {
                        $(".dropzone").animate({"minHeight": 150});
                        $("#secondStep").show();
                        let $canvasArea = $('#canvasArea');
                        $canvasArea.empty();
                        let i = 1;
                        let $pageSelect = $("#page");
                        for(let image of responseText.images) {
                            let id = `canvas_${i}`;
                            let canvas = pdfEditor.createCanvasElement(id);
                            $canvasArea.append(canvas);
                            $pageSelect.append(`<option value="${i}">Page ${i}</option>`);
                            pdfEditor.initPdfEditor(canvas, image, i);
                            i++;
                        }
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
                if (canvas.page == page) {
                    resultCanvas = canvas;
                    break;
                }
            }
        }

        // else {
        //     throw "No canvas element";
        // }
        //
        // if (resultCanvas === null) {
        //     throw "No canvas element. Upload File one more time";
        // }

        return resultCanvas;
    },
    addText: function (e) {
        let $newText = $('#newtext');
        let text = new fabric.Text($newText.val(), {
            fontSize: parseInt($('#fontSize-control').val(), 10),
            fontWeight: parseInt($('#fontWeight-control').val(), 10),
            charSpacing: parseInt($('#charSpacing-control').val(), 10),
        });

        let color = $('#color-control').val();

        if (color === "") {
            color = "#000";
        }

        text.setColor(color);

        let canvasData = pdfEditor.getSelectedCanvas();
        let CANVAS = canvasData.canvas;

        CANVAS.add(text);

        CANVAS.on("object:selected", function (e) {
            pdfEditor.canvasListener(CANVAS, e.target);
        });
    },
    addTextByClick: function (canvas, e) {
        const text = new fabric.IText('...', {
            fontSize: parseInt($('#fontSize-control').val(), 10),
            fontWeight: parseInt($('#fontWeight-control').val(), 10),
            charSpacing: parseInt($('#charSpacing-control').val(), 10),
            top:e.offsetY,
            cursorDuration:500,
            left:e.offsetX,
        });

        canvas.add(text);
    },
    canvasListener: function (CANVAS, elem) {
        const $angleControl = $('#angle-control');
        $angleControl.on("input", function () {
            elem.set('angle', parseInt(this.value, 10)).setCoords();
            CANVAS.renderAll();
        });

        const $scaleControl = $('#scale-control');
        $scaleControl.on("input", function () {
            elem.scale(parseFloat(this.value)).setCoords();
            $fontSizeControl.val((elem.fontSize * elem.scaleX).toFixed(0));
            CANVAS.renderAll();
        });

        const $topControl = $('#top-control');
        $topControl.on("input", function () {
            elem.set('top', parseInt(this.value, 10)).setCoords();
            CANVAS.renderAll();
        });

        const $leftControl = $('#left-control');
        $leftControl.on("input", function () {
            elem.set('left', parseInt(this.value, 10)).setCoords();
            CANVAS.renderAll();
        });

        const $skewXControl = $('#skewX-control');
        $skewXControl.on("input", function () {
            elem.set('skewX', parseInt(this.value, 10)).setCoords();
            CANVAS.renderAll();
        });

        const $skewYControl = $('#skewY-control');
        $skewYControl.on("input", function () {
            elem.set('skewY', parseInt(this.value, 10)).setCoords();
            CANVAS.renderAll();
        });

        const $charSpacingControl = $('#charSpacing-control');
        $charSpacingControl.on("input", function () {
            elem.set('charSpacing', parseInt(this.value, 10)).setCoords();
            CANVAS.renderAll();
        });

        const $fontSizeControl = $('#fontSize-control');
        $fontSizeControl.on("change", function () {
            elem.set('fontSize', parseFloat(this.value));
            CANVAS.renderAll();
        });

        const $colorControl = $('#color-control');
        $colorControl.on("change", function () {
            elem.set('fill', '#' + this.value);
            CANVAS.renderAll();
        });

        const $fontWeightControl = $('#fontWeight-control');
        $fontWeightControl.on("change", function () {
            elem.set('fontWeight', parseInt(this.value, 10));
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
            $fontSizeControl.val(elem.fontSize);
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
                objectsInGroup.forEach(function(object) {
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
            "url": "/api/v1/pdf/save",
            "data": JSON.stringify(results)
        }).done(function (result) {

        }).fail(function () {

        });
    }
};