//
//  ContentView.swift
//  Conversa
//
//  Created by zy on 26/07/22.
//

import SwiftUI
import CoreData

struct ContentView: View {
    init(){
        UITextView.appearance().backgroundColor = .clear
        UITableView.appearance().backgroundColor = .clear
    }
    @Environment(\.managedObjectContext) private var viewContext
    
    @FetchRequest(
        sortDescriptors: [NSSortDescriptor(keyPath: \Item.timestamp, ascending: true)],
        animation: .default)
    private var items: FetchedResults<Item>
    
    @State private var isActive: Bool = false
    @State private var searchText = ""
    @State private var titleText = ""
    @State private var descriptionText = ""
    @FocusState private var isFocused: Bool
    @State var createQuestion: Bool = false
    
    var body: some View {
        NavigationView {
            ZStack {
                LinearGradient(gradient: Gradient(colors: [Color(hex: "#FFFCF4"), Color(hex: "#FFF4D9"), Color(hex: "#FFF7E1")]), startPoint: .top, endPoint: .bottom)
                    .edgesIgnoringSafeArea(.vertical)
                VStack {
                    SearchBar(text: $searchText)
                        .frame(
                            minWidth: 0,
                            maxWidth: 350
                        ).padding(.top, 20)
                    ScrollView {
                        ForEach(items.filter({ searchText.isEmpty ? true : $0.title!.contains(searchText)})) { item in
                            NavigationLink(destination: DetailQuestion()){
                                VStack(alignment: .leading, spacing: 6) {
                                    Text(item.title ?? "Error Title!")
                                        .font(.subheadline)
                                        .bold()
                                        .foregroundColor(Color(hex: "#383838"))
                                        .padding(.horizontal,24)
                                        .padding(.top, 16)
                                        .multilineTextAlignment(.leading)
                                    Text(item.content ?? "Error Content!")
                                        .font(.caption)
                                        .foregroundColor(Color(hex: "#484848"))
                                        .padding(.horizontal,24)
                                        .padding(.top, 4)
                                        .multilineTextAlignment(.leading)
                                }
                                .frame(
                                    minWidth: 0,
                                    maxWidth: 327,
                                    minHeight: 0,
                                    maxHeight: 124,
                                    alignment: .leading
                                )
                                .padding(.bottom)
                                .background(Color.white)
                                .cornerRadius(20)
                                .shadow(color: Color(hex: "#000000").opacity(0.05), radius: 3, x: 0, y: 3)
                                .clipped()
                            }
                            .navigationBarTitleDisplayMode(.inline)
                            // this sets the Back button text when a new screen is pushed
                            .navigationTitle("Questions")
                        }
                        .onDelete(perform: deleteItems)
                        .frame(maxWidth: .infinity)
                    }
                    .padding(.top, 20)
                    .toolbar {
                        ToolbarItem(placement: .principal) {
                            // this sets the screen title in the navigation bar, when the screen is visible
                            Text("")
                        }
                        ToolbarItem(placement: .navigationBarLeading) {
                            Text("Questions")
                                .font(.title2)
                                .bold()
                                .foregroundColor(Color(hex: "#383838"))
                                .padding(.horizontal,16)
                                .multilineTextAlignment(.leading)
                        }
                        ToolbarItem(placement: .navigationBarTrailing) {
                            EditButton()
                        }
                        ToolbarItem {
                            Button(action: {createQuestion = true}) {
                                Label("Add Item", systemImage: "plus")
                            }
                        }
                    }
                }.sheet(isPresented: $createQuestion) {
                    CreateQuestionModal(
                        isPresented: $createQuestion,
                        title: $titleText,
                        description: $descriptionText)
                }
            }
        }
    }
    
    private func addItem() {
        withAnimation {
            let newItem = Item(context: viewContext)
            newItem.timestamp = Date()
            newItem.title = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
            newItem.content = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
            
            do {
                try viewContext.save()
            } catch {
                // Replace this implementation with code to handle the error appropriately.
                // fatalError() causes the application to generate a crash log and terminate. You should not use this function in a shipping application, although it may be useful during development.
                let nsError = error as NSError
                fatalError("Unresolved error \(nsError), \(nsError.userInfo)")
            }
        }
    }
    
    private func deleteItems(offsets: IndexSet) {
        withAnimation {
            offsets.map { items[$0] }.forEach(viewContext.delete)
            
            do {
                try viewContext.save()
            } catch {
                // Replace this implementation with code to handle the error appropriately.
                // fatalError() causes the application to generate a crash log and terminate. You should not use this function in a shipping application, although it may be useful during development.
                let nsError = error as NSError
                fatalError("Unresolved error \(nsError), \(nsError.userInfo)")
            }
        }
    }
}

private let itemFormatter: DateFormatter = {
    let formatter = DateFormatter()
    formatter.dateStyle = .short
    formatter.timeStyle = .medium
    return formatter
}()

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView().environment(\.managedObjectContext, PersistenceController.preview.container.viewContext)
    }
}
