
$(document).on("page-create page-update", function(){
    doSemanticUI("body > #main")
})

var map_individus;
var gmarkers = [];

function mapIGN(map) {
    L.geoportalLayer.WMTS({
        layer: "ORTHOIMAGERY.ORTHOPHOTOS"
    }).addTo(map);    
    L.geoportalLayer.WMTS({
        layer: "GEOGRAPHICALGRIDSYSTEMS.MAPS"
    }).addTo(map);

    var layerSwitcher = L.geoportalControl.LayerSwitcher();
    map.addControl(layerSwitcher);

    var mp = L.geoportalControl.MousePosition();
    map.addControl(mp);
}

function showMarkers(tag, markers) {
    console.log("show markers")

    $.each(gmarkers, function (i, marker) {
        marker.remove()
    })

    gmarkers = []
    $.each(markers, function (i, location) {
        var position = location.Location
        var marker = L.marker(position, {
            title: location.Infos.join("\n")
        })
        gmarkers.push(marker)
        marker.addTo(map_individus)
    })
}

function createMap(tag, center, zoom) {
    if ($(tag).length == 0) return;

    map_individus = L.map($(tag)[0]).setView([center.lat, center.lng], zoom);
    mapIGN(map_individus)

    map_individus.on("moveend zoomend", function(ev){
        var center = map_individus.getCenter()
        var data = { lat: center.lat, lng: center.lng, zoom: map_individus.getZoom() }
        trigger(null, "map-changed", "map", "page", data)
    })

    console.log("createMap")
}

function refreshMap(center, zoom) {
    console.log("refreshMap")
//    map_individus.setView(center, zoom)
    map_individus.invalidateSize()
}


var editMap;
var editMarker;
function updateEditMap(latitude, longitude) {
    var pos = { lat: latitude, lng: longitude }
    editMarker.setPosition(pos)
    editMap.setCenter(pos)
}

function createEditMap(tag) {
    var latitude = parseFloat($("[name=latitude]").val())
    var longitude = parseFloat($("[name=longitude]").val())
    var position = { lat: 47.626951, lng: 6.997541 }

    if (!isNaN(longitude) && !isNaN(latitude)) {
        position = { lat: latitude, lng: longitude };
    }
    editMap = L.map($(tag)[0], {
        scrollWheelZoom: false
    }).setView(position, 12);
    mapIGN(editMap);

    editMarker = L.marker(position, {
        // title: position.toString(),
        draggable: true
    });
    editMarker.addTo(editMap);
    console.log("createEditMap")

    editMarker.on('dragend', function (event) {
        var position = editMarker.getLatLng()
        // editMarker.setTitle(position.toString())
        trigger(null, "position", "map", "page", position)
    })
}

var previewMap
function createPreviewMap(tag, longitude, latitude) {
    var position = { lat: latitude, lng: longitude };

    previewMap = L.map($(tag)[0], {
        scrollWheelZoom: false 
    }).setView(position, 8);
    mapIGN(previewMap);

    var marker = L.marker(position, {})
    marker.addTo(previewMap)
    console.log("createPreviewMap")
}

function removePreviewMap(tag) {
    previewMap.remove()
}

function doSemanticUI(tag) {
    $(tag).find('.ui.modal').modal({ closable: false }).modal("show")
    $(tag).find('.ui.checkbox').checkbox()
    $(tag).find('.ui.dropdown').dropdown({
        forceSelection: false,
        fullTextSearch: true,
        message: {
            addResult: "Ajouter <b>{term}</b>",
            count: "{count} selectionné(s)",
            maxSelections: 'Max {maxCount} selections',
            noResults: 'Aucun résultat.'
        }
    })
    $(tag).find('.ui.dropdown.additions').dropdown({
        forceSelection: false,
        allowAdditions: true,
        hideAdditions: false,
        message: {
            addResult: 'Ajouter <b>{term}</b>',
            count: '{count} selectionné(s)',
            maxSelections: 'Max {maxCount} selections',
            noResults: 'Aucun résultat.'
        }
    })

    $(tag).find('select.dropdown').dropdown()

}
